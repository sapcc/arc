package db

// Global
var DeleteQuery = "DELETE FROM "
var CheckConnection = "SELECT 1"

// Jobs
var InsertJobQuery = `INSERT INTO jobs(id,version,sender,"to",timeout,agent,action,payload,status,created_at,updated_at,project,"user") VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) returning id;`
var UpdateJobQuery = `UPDATE jobs SET status=$1,updated_at=$2 WHERE id=$3`
var GetAllJobsQuery = "SELECT * FROM jobs %s order by created_at DESC %s"
var CountAllJobsQuery = "SELECT count(*) FROM jobs %s %s"
var GetJobQuery = "SELECT * FROM jobs WHERE id=$1"
var FailJobsTimeoutQuery = `
	UPDATE jobs SET status=3,updated_at=NOW()
	WHERE id IN
	(
		SELECT DISTINCT id
		FROM jobs
		WHERE (created_at <= NOW() - INTERVAL '1 second' * $1 - INTERVAL '1 second' * timeout)
		AND (status=1 OR status=2)
	)
`
var FailJobsNonHeartbeatQuery = `
	UPDATE jobs SET status=3,updated_at=NOW()
	WHERE id IN
	(
		SELECT DISTINCT id
		FROM jobs
		WHERE (created_at <= NOW() - INTERVAL '1 second' * $1)
		AND status=1
	)
`
var DeleteJobsOldQuery = `
	DELETE FROM jobs
	WHERE id IN
	(
		SELECT DISTINCT id
		FROM jobs
		WHERE (updated_at <= NOW() - INTERVAL '1 day' * $1)
		AND status=4
	)
`

// Log
var GetLogQuery = "SELECT * FROM logs WHERE job_id=$1"
var InsertLogQuery = "INSERT INTO logs(job_id,content,created_at,updated_at) VALUES($1,$2,$3,$4) returning job_id"

// Log parts
var GetLogPartQuery = `SELECT * FROM log_parts WHERE job_id=$1 AND number=$2`
var InsertLogPartQuery = `INSERT INTO log_parts(job_id,number,content,final,created_at) VALUES($1,$2,$3,$4,$5) returning job_id;`
var CollectLogPartsQuery = "SELECT array_to_string(array_agg(log_parts.content ORDER BY number, job_id), '') AS content FROM log_parts WHERE job_id=$1"
var DeleteLogPartsQuery = `DELETE FROM log_parts WHERE job_id=$1`
var GetLogPartsToAggregateQuery = `
	SELECT DISTINCT job_id
	FROM log_parts
	WHERE (created_at <= NOW() - INTERVAL '1 seconds' * $1 AND final = true)
	OR created_at <= NOW() - INTERVAL '1 seconds' * $2
`

// Agents
var GetAgentsQuery = `SELECT DISTINCT
CASE WHEN tags->>'name' <> '' THEN tags->>'name' ELSE (
CASE WHEN facts->>'hostname' <> '' THEN facts->>'hostname' ELSE agent_id END
) END as display_name,
*
FROM agents %s order by display_name ASC %s`
var CountAgentsQuery = "SELECT count(*) FROM agents %s %s"
var GetAgentQuery = "SELECT * FROM agents WHERE agent_id=$1"
var InsertAgentQuery = `INSERT INTO agents(agent_id,project,organization,facts,created_at,updated_at,updated_with,updated_by,tags) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) returning agent_id`
var UpdateAgentWithRegistration = `
	UPDATE agents SET
	project=$2,
	organization=$3,
	facts=json_replace((SELECT facts::json FROM agents WHERE agent_id=$1),$4::json)::jsonb,
	updated_at=$5,
	updated_with=$6,
	updated_by=$7
	WHERE agent_id=$1
`
var AddAgentTag = `
	UPDATE agents SET
	updated_at=$2,
	tags=json_set_key((SELECT tags::json FROM agents WHERE agent_id=$1),$3, $4::TEXT)::jsonb
	WHERE agent_id=$1
`
var DeleteAgentQuery = `DELETE FROM agents WHERE agent_id=$1`
var DeleteAgentTagQuery = `
	UPDATE agents SET
	updated_at=$2,
	tags=json_delete_keys((SELECT tags::json FROM agents WHERE agent_id=$1),$3)::jsonb
	WHERE agent_id=$1`

// Locks
var GetLockQuery = "SELECT * FROM locks WHERE lock_id=$1"
var InsertLockQuery = `INSERT INTO locks(lock_id,agent_id,created_at) VALUES($1,$2,$3) returning lock_id`
var DeleteLocksQuery = `
	DELETE FROM locks
	WHERE lock_id IN
	(
		SELECT DISTINCT lock_id
		FROM locks
		WHERE (created_at <= NOW() - INTERVAL '1 seconds' * $1)
	)
`

// pki
var DeletePkiTokensQuery = `DELETE FROM tokens WHERE created_at < NOW() - INTERVAL '1 second' * $1`
var InsertTokenQuery = `INSERT INTO tokens (id, profile, subject) VALUES($1, $2, $3)`
var GetTokenQuery = `SELECT profile, subject FROM tokens WHERE id=$1 AND created_at > NOW() - INTERVAL '1 hour' FOR UPDATE`
var InsertTokenWithCreatedAtQuery = `INSERT INTO tokens (id, profile, subject,created_at) VALUES($1, $2, $3, $4)`
var CleanPkiCertificatesQuery = `DELETE FROM certificates WHERE not_after < NOW()`
var InsertCertificateQuery = `INSERT into certificates
(fingerprint, common_name, country, locality, organization, organizational_unit, not_before, not_after, pem)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`
