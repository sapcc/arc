package db

// Jobs
var InsertJobQuery = `INSERT INTO jobs(id,version,sender,"to",timeout,agent,action,payload,status,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning id;`
var UpdateJobQuery = `UPDATE jobs SET status=$1,updated_at=$2 WHERE id=$3`
var GetAllJobsQuery = "SELECT * FROM jobs order by updated_at"
var GetJobQuery = "SELECT * FROM jobs WHERE id=$1"
var CleanJobsTimeoutQuery = `
	UPDATE jobs SET status=3,updated_at=NOW() 
	WHERE id IN 
	(
		SELECT DISTINCT id
		FROM jobs
		WHERE (created_at <= NOW() - INTERVAL '1 second' * $1 - INTERVAL '1 second' * timeout)
		AND (status=1 OR status=2)
	)
`
var CleanJobsNonHeartbeatQuery = `
	UPDATE jobs SET status=3,updated_at=NOW() 
	WHERE id IN 
	(
		SELECT DISTINCT id
		FROM jobs
		WHERE (created_at <= NOW() - INTERVAL '1 second' * $1)
		AND status=1
	)
`

// Global
var DeleteQuery = "DELETE FROM "

// Log
var GetLogQuery = "SELECT * FROM logs WHERE job_id=$1"
var InsertLogQuery = "INSERT INTO logs(job_id,content,created_at,updated_at) VALUES($1,$2,$3,$4) returning job_id"
var UpdateLogQuery = "UPDATE logs SET content=$1,updated_at=$2 WHERE job_id=$3"

// Log parts
var GetLogPartQuery = `SELECT * FROM log_parts WHERE job_id=$1`
var InsertLogPartQuery = `INSERT INTO log_parts(job_id,number,content,final,created_at) VALUES($1,$2,$3,$4,$5) returning job_id;`
var CollectLogPartsQuery = "SELECT array_to_string(array_agg(log_parts.content ORDER BY number, job_id), '') AS content FROM log_parts WHERE job_id=$1"
var DeleteLogPartsQuery = `DELETE FROM log_parts WHERE job_id=$1`

// Agents
var GetAgentsQuery = "SELECT DISTINCT * FROM agents order by updated_at"
var GetAgentQuery = "SELECT * FROM agents WHERE agent_id=$1"
var InsertAgentQuery = `INSERT INTO agents(agent_id,project,organization,facts,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6) returning agent_id`
var UpdateAgent = `UPDATE agents SET project=$2,organization=$3,facts=json_replace((SELECT facts::json FROM agents WHERE agent_id=$1),$4::json)::jsonb,updated_at=$5 WHERE agent_id=$1`
