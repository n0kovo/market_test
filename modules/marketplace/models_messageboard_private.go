package marketplace

// func setupPrivateThreadsFunctions() {
// 	database.Exec("DROP VIEW IF EXISTS v_private_threads CASCADE;")
// 	database.Exec(`
// 		CREATE VIEW v_private_threads AS (
// 			SELECT
// 				v_threads.*,
// 				('{' || v_threads.sender_user_uuid || ', ' || v_threads.reciever_user_uuid || '}')::text[] as participants
// 			FROM v_threads
// 			WHERE type='private'
// 	);`)
// }
