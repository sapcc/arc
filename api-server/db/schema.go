package db

// Jobs
var InsertJobQuery = `INSERT INTO jobs(version,sender,requestid,"to",timeout,agent,action,payload,status,createdat,updatedat) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning requestid;`
var UpdateJobQuery = `UPDATE jobs SET status=$1,updatedat=$2 WHERE requestid=$3`
var GetAllJobsQuery = "SELECT * FROM jobs order by requestid"
var GetJobQuery = "SELECT * FROM jobs WHERE requestid=$1"
var CleanJobsTimeoutQuery = `
	UPDATE jobs SET status=3,updatedat=NOW() 
	WHERE requestid IN 
	(
		SELECT DISTINCT requestid
		FROM jobs
		WHERE (createdat <= NOW() - INTERVAL '1 second' * $1 - INTERVAL '1 second' * timeout)
		AND (status=1 OR status=2)
	)
`
var CleanJobsNonHeartbeatQuery = `
	UPDATE jobs SET status=3,updatedat=NOW() 
	WHERE requestid IN 
	(
		SELECT DISTINCT requestid
		FROM jobs
		WHERE (createdat <= NOW() - INTERVAL '1 second' * $1)
		AND status=1
	)
`

// Log
var GetLogQuery = "SELECT * FROM logs WHERE job_id=$1"
var InsertLogQuery = "INSERT INTO logs(job_id,content,createdat,updatedat) VALUES($1,$2,$3,$4) returning job_id"
var UpdateLogQuery = "UPDATE logs SET content=$1,updatedat=$2 WHERE job_id=$3"

// Log parts
var InsertLogPartQuery = `INSERT INTO log_parts(job_id,number,content,final,createdat) VALUES($1,$2,$3,$4,$5) returning job_id;`
var CollectLogPartsQuery = "SELECT array_to_string(array_agg(log_parts.content ORDER BY number, job_id), '') AS content FROM log_parts WHERE job_id=$1"
var DeleteLogPartsQuery = `DELETE FROM log_parts WHERE job_id=$1`

// Facts
var GetAgentsQuery = "SELECT DISTINCT agent_id,createdat,updatedat FROM facts"
var GetAgentQuery = "SELECT agent_id,createdat,updatedat FROM facts WHERE agent_id=$1"
var GetFactQuery = "SELECT * FROM facts WHERE agent_id=$1"
var InsertFactQuery = `INSERT INTO facts(agent_id,facts,createdat,updatedat) VALUES($1,$2,$3,$4) returning agent_id`
var UpdateFact = `UPDATE facts SET facts=json_replace((SELECT facts::json FROM facts WHERE agent_id=$1),$2::json)::jsonb where agent_id=$1`
