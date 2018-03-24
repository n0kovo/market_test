package marketplace

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/gocraft/web"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) UserAvatar(w web.ResponseWriter, r *web.Request) {
	size := "normal"
	if len(r.URL.Query()["size"]) > 0 {
		size = r.URL.Query()["size"][0]
	}

	username := r.PathParams["user"]
	userUuid := CacheGetUserUuid(username)
	if userUuid != "" {
		err := util.ServeImage(userUuid+"_av", size, w, r)
		if err != nil {
			user, _ := FindUserByUsername(username)
			user.GenerateAvatar()

			util.ServeImage(userUuid+"_av", size, w, r)
		}
	} else {
		http.NotFound(w, r.Request)
		return
	}
}

func (c *Context) Login(user User, w web.ResponseWriter, r *web.Request) {

	if c.ViewUser.Uuid != "" && !c.ViewUser.IsAdmin {
		http.NotFound(w, r.Request)
		return
	}

	session, _ := sessionStore.Get(r.Request, "auth-session")
	session.Values["UserUuid"] = user.Uuid
	session.Save(r.Request, w)

	now := time.Now()
	user.LastLoginDate = &now
	user.Save()

	redirectUrl := "/marketplace"
	EventUserLoggedIn(user)
	http.Redirect(w, r.Request, redirectUrl, 302)
}

func (c *Context) RegisterGET(w web.ResponseWriter, r *web.Request) {

	if c.ViewUser.Uuid != "" {
		http.NotFound(w, r.Request)
		return
	}

	if r.PathParams["invite_code"] != "" {
		c.InviteCode = r.PathParams["invite_code"]
		user, _ := FindUserByInviteCode(r.PathParams["invite_code"])
		if user != nil {
			seller := Seller{user}
			c.ViewSeller = seller.ViewSeller(c.ViewUser.User.Language)
		} else {
			c.Error = "Invite code is invalid"
		}
	}

	c.SelectedSection = "register"
	c.CaptchaId = captcha.New()
	util.RenderTemplate(w, "auth/register", c)
}

func (c *Context) ViewRecoverGET(w web.ResponseWriter, r *web.Request) {
	if c.ViewUser.Uuid != "" {
		http.NotFound(w, r.Request)
		return
	}

	c.CaptchaId = captcha.New()

	util.RenderTemplate(w, "auth/recovery_step_1", c)
}

func (c *Context) ViewRecoverPOST(w web.ResponseWriter, r *web.Request) {

	c.CaptchaId = captcha.New()

	var (
		session, _ = sessionStore.Get(r.Request, "auth-session")
	)

	setEncryptedText := func(user *User) {
		secretText := util.GenerateUuid()

		session.Values["2FactorUserUuid"] = user.Uuid
		session.Values["secrettext"] = secretText
		session.Save(r.Request, w)

		c.SecretText, _ = util.EncryptText(secretText, user.Pgp)
	}

	switch r.FormValue("step") {
	case "1":
		isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
		if !isCaptchaValid {
			c.Error = "Invalid captcha"
			util.RenderTemplate(w, "auth/recovery_step_1", c)
			return
		}

		user, err := FindUserByUsername(r.FormValue("username"))
		if err != nil {
			c.Error = "User not found"
			util.RenderTemplate(w, "auth/recovery_step_1", c)
			break
		}

		if user.Pgp == "" {
			c.Error = "User doesn't have PGP set up"
			util.RenderTemplate(w, "auth/recovery_step_1", c)
			break
		}

		setEncryptedText(user)
		util.RenderTemplate(w, "auth/recovery_step_2", c)
		break

	case "2":
		var (
			secretText, _    = (session.Values["secrettext"]).(string)
			userId, _        = (session.Values["2FactorUserUuid"]).(string)
			decryptedmessage = strings.Trim(r.FormValue("secret_text"), "\n ")
		)

		user, _ := FindUserByUuid(userId, false)
		if user == nil {
			c.Error = "Could not find user"
			util.RenderTemplate(w, "auth/recovery_step_2", c)
			return
		}

		if decryptedmessage == secretText {
			session.Values["UserUuid"] = user.Uuid
			session.Save(r.Request, w)
			util.RenderTemplate(w, "auth/recovery_step_3", c)
			return
		} else {
			c.Error = "Could not authenticate"
			setEncryptedText(user)
			util.RenderTemplate(w, "auth/recovery_step_2", c)
			return
		}
		break

	case "3":
		var (
			userId, _ = (session.Values["UserUuid"]).(string)
			user, _   = FindUserByUuid(userId, false)
		)

		if user == nil || userId == "" {
			c.Error = "Could not find user"
			util.RenderTemplate(w, "auth/recovery_step_3", c)
			return
		}

		if r.FormValue("passphrase") != r.FormValue("passphrase_2") {
			c.Error = "Passphrases do not match"
			util.RenderTemplate(w, "auth/recovery_step_3", c)
			return
		}

		user.PassphraseHash = util.PasswordHashV1(user.Username, r.FormValue("passphrase"))
		user.Save()

		http.Redirect(w, r.Request, "/marketplace", 302)
		break

	default:
		http.NotFound(w, r.Request)
		return
	}

}

func (c *Context) RegisterPOST(w web.ResponseWriter, r *web.Request) {
	if c.ViewUser.Uuid != "" {
		http.NotFound(w, r.Request)
		return
	}
	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.RegisterGET(w, r)
		return
	}
	if r.FormValue("passphrase") != r.FormValue("passphrase_2") {
		c.Error = "Passphrases do not match"
		c.RegisterGET(w, r)
		return
	}

	var inviterUuid string
	if r.FormValue("invite_code") != "" {
		inviter, _ := FindUserByInviteCode(r.FormValue("invite_code"))
		if inviter != nil {
			inviterUuid = inviter.Uuid
		} else {
			c.Error = "Invite code is invalid"
			c.RegisterGET(w, r)
			return
		}
	}

	user, validationError := CreateUser(r.FormValue("username"), r.FormValue("passphrase"))
	if validationError != nil {
		c.Error = validationError.Error()
		c.RegisterGET(w, r)
		return
	}

	if user.Language == "" {
		lang := "en"
		if c.Language != "" {
			lang = c.Language
		}
		user.Language = lang
	}

	session, _ := sessionStore.Get(r.Request, "auth-session")
	session.Values["UserUuid"] = user.Uuid
	session.Save(r.Request, w)

	now := time.Now()
	user.LastLoginDate = &now
	if inviterUuid != "" {
		user.InviterUuid = inviterUuid
	}

	if r.FormValue("role") == "seller" {
		user.IsSeller = true

	}
	user.Save()

	var redirectUrl string
	if user.IsSeller {
		redirectUrl = "/profile"
	} else {
		redirectUrl = "/marketplace"
	}
	EventUserRegistred(*user)
	http.Redirect(w, r.Request, redirectUrl, 302)
}

func (c *Context) BotCheckGET(w web.ResponseWriter, r *web.Request) {
	c.CaptchaId = captcha.New()
	util.RenderTemplate(w, "auth/bot-check", c)
}

func (c *Context) BotCheckPOST(w web.ResponseWriter, r *web.Request) {

	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.RegisterGET(w, r)
		return
	}

	session, _ := sessionStore.Get(r.Request, "auth-session")
	uuid := util.GenerateUuid()
	session.Values["BotCheckUuid"] = uuid
	botCheckUuids[uuid] = true
	session.Save(r.Request, w)

	if redirect := r.Request.URL.Query().Get("redirect"); redirect != "" {
		http.Redirect(w, r.Request, redirect, 302)
	} else {
		http.Redirect(w, r.Request, "/", 302)
	}
}

func (c *Context) LoginGET(w web.ResponseWriter, r *web.Request) {

	if c.ViewUser.Uuid != "" {
		http.NotFound(w, r.Request)
		return
	}
	c.SelectedSection = "login"

	c.CaptchaId = captcha.New()
	util.RenderTemplate(w, "auth/login", c)
}

func (c *Context) AccountGET(w web.ResponseWriter, r *web.Request) {
	c.ViewSeller = Seller{c.ViewUser.User}.ViewSeller(c.ViewUser.User.Language) //@
	c.USDBTCRate = GetCurrencyRate("BTC", "USD")
	util.RenderTemplate(w, "auth/account", c)
}

func (c *Context) AccountPOST(w web.ResponseWriter, r *web.Request) {

	err := r.ParseForm()
	if err != nil {
		c.Error = err.Error()
		c.AccountGET(w, r)
		return
	}

	paymentType := r.FormValue("type")

	if c.ViewUser.User.PremiumPlus {
		c.Error = "You are alredy Premuim+ vendor"
		c.AccountGET(w, r)
		return
	}

	if c.ViewUser.User.Premium && paymentType == "premium" {
		c.Error = "You are alredy Premuim vendor"
		c.AccountGET(w, r)
		return
	}

	var priceUSD float64
	c.USDBTCRate = GetCurrencyRate("BTC", "USD")

	switch paymentType {
	case "premium":
		priceUSD = 50
	case "premium_plus":
		if c.ViewUser.User.Premium {
			priceUSD = 50
		} else {
			priceUSD = 100
		}
	}

	price := priceUSD / c.USDBTCRate

	userWallets := c.ViewUser.User.FindUserBitcoinWallets()
	if userWallets.Balance().Balance < price {
		c.Error = fmt.Sprintf("Please deposit %f BTC to your onsite wallet.", price)
		c.AccountGET(w, r)
		return
	}

	addr, err := apis.GenerateBTCAddress("premium")
	if err != nil {
		c.Error = err.Error()
		c.AccountGET(w, r)
		return
	}

	_, err = userWallets.Send(addr, price)
	if err != nil {
		c.Error = err.Error()
		c.AccountGET(w, r)
		return
	}

	switch paymentType {
	case "premium":
		c.ViewUser.User.Premium = true
	case "premium_plus":
		c.ViewUser.User.Premium = true
		c.ViewUser.User.PremiumPlus = true
	}

	c.ViewUser.User.Save()
	c.AccountGET(w, r)
}

func (c *Context) getUserForTrustPage(r *web.Request) (User, error) {
	user := *c.ViewUser.User
	if r.PathParams["username"] != "" && !(c.ViewUser.IsStaff || c.ViewUser.IsAdmin) {
		return User{}, errors.New("Wrong request")
	}
	if r.PathParams["username"] != "" && (c.ViewUser.IsStaff || c.ViewUser.IsAdmin) {
		u, err := FindUserByUsername(r.PathParams["username"])
		if err != nil {
			return User{}, errors.New("Wrong request")
		}
		user = *u
	}
	return user, nil
}

func (c *Context) ViewTrustGET(w web.ResponseWriter, r *web.Request) {

	user, err := c.getUserForTrustPage(r)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	c.ViewSeller = Seller{&user}.ViewSeller(c.ViewUser.Language)

	thread, err := GetVendorVerificationThread(user, true)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	c.CaptchaId = captcha.New()
	c.ViewThread = thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
	util.RenderTemplate(w, "auth/trust", c)
}

func (c *Context) ViewTrustPOST(w web.ResponseWriter, r *web.Request) {

	if r.FormValue("request_trusted_vendor") != "" {
		c.ViewUser.HasRequestedVerification = true
		c.ViewUser.User.Save()
		c.ViewTrustGET(w, r)
		return
	}

	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ViewTrustGET(w, r)
		return
	}

	user, err := c.getUserForTrustPage(r)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	thread, err := GetVendorVerificationThread(user, false)
	if err != nil {
		c.Error = err.Error()
		c.ViewTrustGET(w, r)
		return
	}
	message, err := CreateMessage(r.FormValue("text"), *thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewTrustGET(w, r)
		return
	}

	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ViewTrustGET(w, r)
		return
	}

	EventNewTrustedVendorThreadPost(*c.ViewUser.User, Seller{&user}, *message)
	c.ViewTrustGET(w, r)
}

func (c *Context) Referrals(w web.ResponseWriter, r *web.Request) {
	if c.ViewUser.InviteCode == "" {
		c.ViewUser.User.InviteCode = util.GenerateUuid()
		c.ViewUser.User.Save()
	}

	c.SelectedSection = "users"
	if len(r.URL.Query()["section"]) > 0 {
		c.SelectedSection = r.URL.Query()["section"][0]
	}

	// paging
	pageSize := 50
	selectedPage := 0
	if len(r.URL.Query()["page"]) > 0 {
		selectedPageStr := r.URL.Query()["page"][0]
		page, err := strconv.Atoi(selectedPageStr)
		if err == nil {
			selectedPage = page - 1
		}
	}
	numberOfPages := int(math.Ceil(float64(c.NumberOfInvitedUsers) / float64(pageSize)))
	for i := 0; i < numberOfPages; i++ {
		c.Pages = append(c.Pages, i+1)
	}

	usersPage := GetInvitedUserPage(c.ViewUser.Uuid, selectedPage, pageSize)
	c.ExtendedUsers = usersPage
	c.ReferralPayments = FindReferralPaymentsForUser(c.ViewUser.User.Uuid)
	c.NumberOfInvitedUsers = c.ViewUser.NumberOfInvitedUsers()

	util.RenderTemplate(w, "auth/referrals", c)
}

func (c *Context) Banned(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplate(w, "auth/banned", c)
}

func (c *Context) FreeRestrictions(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplate(w, "free_restrictions", c)
}

func (c *Context) Logout(w web.ResponseWriter, r *web.Request) {
	session, _ := sessionStore.Get(r.Request, "auth-session")
	delete(session.Values, "UserUuid")
	session.Save(r.Request, w)
	http.Redirect(w, r.Request, "/", 302)
}

func (c *Context) LoginPOST(w web.ResponseWriter, r *web.Request) {
	if r.FormValue("decryptedmessage") == "" {
		c.Login1FactorPOST(w, r)
	} else {
		c.Login2FactorPOST(w, r)
	}
}

func (c *Context) Login1FactorPOST(w web.ResponseWriter, r *web.Request) {
	isCaptchaValid := captcha.VerifyString(r.FormValue("captcha_id"), r.FormValue("captcha"))
	user, _ := FindUserByUsername(r.FormValue("username"))
	isLoginSuccessful := isCaptchaValid && (user != nil) && user.CheckPassphrase(r.FormValue("passphrase"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.LoginGET(w, r)
		return
	}
	if user == nil {
		c.Error = "Failed to authenticate"
		c.LoginGET(w, r)
		return
	}
	if !isLoginSuccessful {
		c.Error = "Failed to authenticate"
		c.LoginGET(w, r)
		return
	}
	if user.TwoFactorAuthentication {
		session, _ := sessionStore.Get(r.Request, "auth-session")
		session.Values["2FactorUserUuid"] = user.Uuid
		// session.Values["2fasecret"] = uuid
		session.Save(r.Request, w)
		c.LoginPGPGet(w, r)
	} else {
		c.Login(*user, w, r)
	}
}

func (c *Context) Login2FactorPOST(w web.ResponseWriter, r *web.Request) {
	session, _ := sessionStore.Get(r.Request, "auth-session")

	secretText, _ := (session.Values["secrettext"]).(string)
	userId, _ := (session.Values["2FactorUserUuid"]).(string)
	decryptedmessage := strings.Trim(r.FormValue("decryptedmessage"), "\n ")

	user, _ := FindUserByUuid(userId, false)
	if user == nil {
		c.Error = "Could not authenticate"
		c.LoginPGPGet(w, r)
		return
	}

	isSingatureCorrect := decryptedmessage == secretText
	if isSingatureCorrect {
		c.Login(*user, w, r)
		return
	} else {
		c.Error = "Could not authenticate"
		c.LoginPGPGet(w, r)
		return
	}
}

func (c *Context) LoginPGPGet(w web.ResponseWriter, r *web.Request) {

	session, _ := sessionStore.Get(r.Request, "auth-session")
	userId, _ := (session.Values["2FactorUserUuid"]).(string)

	user, _ := FindUserByUuid(userId, false)
	if user == nil {
		c.Error = "Could not authenticate"
		c.LoginGET(w, r)
		return
	}

	secretText := util.GenerateUuid()

	session.Values["secrettext"] = secretText
	session.Save(r.Request, w)

	c.SecretText, _ = util.EncryptText(secretText, user.Pgp)
	util.RenderTemplate(w, "auth/pgplogin", c)
}

func (c *Context) UserBanner(w web.ResponseWriter, r *web.Request) {
	user, _ := FindUserByUsername(r.PathParams["user"])
	if user == nil || !user.HasTopBanner {
		http.NotFound(w, r.Request)
		return
	}
	size := "728x90"
	util.ServeImage(user.Uuid+"_tb", size, w, r)
}

func (c *Context) ProfileGET(w web.ResponseWriter, r *web.Request) {
	if len(r.URL.Query()["section"]) > 0 {
		section := r.URL.Query()["section"][0]
		c.SelectedSection = section
	} else {
		c.SelectedSection = "profile"
	}
	secretText := util.GenerateUuid()
	session, _ := sessionStore.Get(r.Request, "auth-session")
	session.Values["secrettext"] = secretText
	session.Save(r.Request, w)
	c.SecretText = secretText

	if c.ViewUser.IsSeller {
		c.ViewUser.User, _ = FindUserByUuid(c.ViewUser.Uuid, true)
	}

	c.UserSettingsHistory = SettingsChangeHistoryByUser(c.ViewUser.User.Uuid)

	util.RenderTemplate(w, "auth/profile", c)
}

func (c *Context) ProfilePOST(w web.ResponseWriter, r *web.Request) {

	var (
		previousBTCAddress = c.ViewUser.User.Bitcoin
		previousBCHAddress = c.ViewUser.User.BitcoinCash
		previousETHAddress = c.ViewUser.User.Ethereum
		btcAddress         = r.FormValue("bitcoin")
		bchAddress         = r.FormValue("bitcoin_cash")
		ethereumAddress    = r.FormValue("ethereum")
	)

	if r.FormValue("description") != "" {
		c.ViewUser.User.Description = r.FormValue("description")
	}
	if r.FormValue("long_description") != "" {
		c.ViewUser.User.LongDescription = r.FormValue("long_description")
	}

	if btcAddress != "" && !bitcoinRegexp.MatchString(btcAddress) {
		c.Error = "Wrong Bitcoin Address"
	} else {
		c.ViewUser.User.Bitcoin = btcAddress
	}

	if bchAddress != "" && !bitcoinRegexp.MatchString(bchAddress) {
		c.Error = "Wrong Bitcoin Cash Address"
	} else {
		c.ViewUser.User.BitcoinCash = bchAddress
	}

	if ethereumAddress != "" && !ethereumRegexp.MatchString(ethereumAddress) {
		c.Error = "Wrong Ethereum"
	} else {
		c.ViewUser.User.Ethereum = ethereumAddress
	}

	if r.FormValue("bitcoin_multisig") != "" {
		c.ViewUser.User.BitcoinMultisigPublicKey = r.FormValue("bitcoin_multisig")
	}
	if r.FormValue("bitmessage") != "" {
		c.ViewUser.User.Bitmessage = r.FormValue("bitmessage")
	}
	if r.FormValue("tox") != "" {
		c.ViewUser.User.Tox = r.FormValue("tox")
	}
	if r.FormValue("email") != "" {
		c.ViewUser.User.Email = r.FormValue("email")
	}
	if r.FormValue("2fa") != "" {
		if r.FormValue("2fa") == "1" {
			c.ViewUser.User.TwoFactorAuthentication = true
		} else if r.FormValue("2fa") == "0" {
			c.ViewUser.User.TwoFactorAuthentication = false
		}
	}
	if r.FormValue("vm") != "" {
		if r.FormValue("vm") == "1" {
			c.ViewUser.User.VacationMode = true
		} else if r.FormValue("vm") == "0" {
			c.ViewUser.User.VacationMode = false
		}
	}
	if r.FormValue("old_password") != "" {
		oldPassword := r.FormValue("old_password")
		hashV1 := util.PasswordHashV1(c.ViewUser.User.Username, oldPassword)
		if c.ViewUser.User.PassphraseHash != hashV1 {
			c.Error = "Invalid worng password"
			c.ProfileGET(w, r)
			return
		}

		newPassword := r.FormValue("new_password")
		repeatNewPassword := r.FormValue("repeat_new_password")

		if newPassword != repeatNewPassword {
			c.Error = "New password and repeat new password does not match"
			c.ProfileGET(w, r)
			return
		}

		newHash := util.PasswordHashV1(c.ViewUser.User.Username, newPassword)
		c.ViewUser.User.PassphraseHash = newHash
		c.ViewUser.User.Save()

	}
	if validationError := c.ViewUser.User.Validate(); validationError != nil {
		c.Error = validationError.Error()
		c.ProfileGET(w, r)
		return
	}
	avatarError := util.SaveImage(r, "avatar_image", 300, c.ViewUser.User.Uuid+"_av")
	if avatarError == nil {
		c.ViewUser.User.HasAvatar = true
	}

	tbError := util.SaveImage(r, "top_banner_image", 728, c.ViewUser.User.Uuid+"_tb")
	if tbError == nil {
		c.ViewUser.User.HasTopBanner = true
	}

	c.ViewUser.User.Save()

	if previousBTCAddress != c.ViewUser.User.Bitcoin {
		historyEvent := UserSettingsHistory{
			UserUuid: c.ViewUser.User.Uuid,
			Action:   "Bitcoin address changed to " + c.ViewUser.User.Bitcoin,
			Datetime: time.Now(),
			Type:     "bitcoin",
		}
		if c.ViewUser.User.Bitcoin == "" {
			historyEvent.Action = "Bitcoin address deleted"
		}
		historyEvent.Save()
	}

	if previousBCHAddress != c.ViewUser.User.BitcoinCash {
		historyEvent := UserSettingsHistory{
			UserUuid: c.ViewUser.User.Uuid,
			Action:   "BitcoinCash address changed to " + c.ViewUser.User.BitcoinCash,
			Datetime: time.Now(),
			Type:     "bitcoin_cash",
		}
		if c.ViewUser.User.BitcoinCash == "" {
			historyEvent.Action = "Bitcoin Cash address deleted"
		}
		historyEvent.Save()
	}

	if previousETHAddress != c.ViewUser.User.Ethereum {
		historyEvent := UserSettingsHistory{
			UserUuid: c.ViewUser.User.Uuid,
			Action:   "Ethereum address changed to " + c.ViewUser.User.Ethereum,
			Datetime: time.Now(),
			Type:     "ethereum",
		}
		if c.ViewUser.User.Ethereum == "" {
			historyEvent.Action = "Ethereum address deleted"
		}
		historyEvent.Save()
	}

	c.ProfileGET(w, r)
}

func (c *Context) SetupPGPViaDecryptionStep1GET(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplate(w, "auth/settings_pgp_signature_step_1", c)
}

func (c *Context) SetupPGPViaDecryptionStep1POST(w web.ResponseWriter, r *web.Request) {
	pgp := r.FormValue("pgp_public_key")
	uuid := util.GenerateUuid()

	encryptedPgp, err := util.EncryptText(uuid, pgp)
	if err != nil {
		c.Error = err.Error()
		c.SetupPGPViaDecryptionStep1GET(w, r)
		return
	}
	c.SecretText = encryptedPgp

	session, _ := sessionStore.Get(r.Request, "auth-session")
	session.Values["uuid"] = uuid
	session.Save(r.Request, w)

	c.Pgp = pgp
	util.RenderTemplate(w, "auth/settings_pgp_signature_step_2", c)
}

func (c *Context) SetupPGPViaDecryptionStep2POST(w web.ResponseWriter, r *web.Request) {
	decryptedText := r.FormValue("secret_text")

	session, _ := sessionStore.Get(r.Request, "auth-session")
	secretText, _ := (session.Values["uuid"]).(string)
	pgp := r.FormValue("pgp")

	if decryptedText == secretText {
		c.ViewUser.User.Pgp = strings.Trim(pgp, "\n ")
		c.ViewUser.User.Save()
		redirectUrl := "/profile?section=system"
		http.Redirect(w, r.Request, redirectUrl, 302)
		return
	} else {
		c.Error = "Decrypted message isn't correct"
		c.SetupPGPViaDecryptionStep1GET(w, r)
		return
	}
}

func (c *Context) SetCurrency(w web.ResponseWriter, r *web.Request) {
	currency := r.PathParams["currency"]

	allowedCurrencies := map[string]bool{
		"AUD": true,
		"BTC": true,
		"ETH": true,
		"EUR": true,
		"GBP": true,
		"RUB": true,
		"USD": true,
	}

	if _, ok := allowedCurrencies[currency]; ok {
		c.ViewUser.User.Currency = currency
		c.ViewUser.User.Save()
	}

	redirectUrl := r.Referer()
	http.Redirect(w, r.Request, redirectUrl, 302)
}

func (c *Context) SetLanguage(w web.ResponseWriter, r *web.Request) {
	language := r.PathParams["language"]

	redirectUrl := r.Referer()
	if strings.Contains(redirectUrl, "/settings/language") {
		redirectUrl = "/"
	}

	if language == "ru" || language == "en" || language == "es" || language == "fr" || language == "de" || language == "rs" || language == "tr" || language == "it" {
		c.ViewUser.User.Language = language
		c.ViewUser.User.Save()
	}

	http.Redirect(w, r.Request, redirectUrl, 302)
}
