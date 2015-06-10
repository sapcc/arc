package db

var jobsTable = `
	CREATE TABLE IF NOT EXISTS jobs
	(
		version integer NOT NULL,
		sender text NOT NULL,
		requestid text PRIMARY KEY NOT NULL,
		"to" text NOT NULL,
		timeout integer NOT NULL,
		agent text NOT NULL,
		action text NOT NULL,
		payload text NOT NULL,
		status text NOT NULL,
		CONSTRAINT uc_requestid UNIQUE (requestid)		
	)
	WITH (
 	 OIDS=FALSE
	);
	ALTER TABLE jobs
		OWNER TO arc;
`

var agentsTable = `
	CREATE TABLE IF NOT EXISTS agents 
	( 
		uid integer PRIMARY KEY NOT null, 
		CONSTRAINT uc_uid UNIQUE (uid)
	)
	WITH (
	 OIDS=FALSE
	);
	ALTER TABLE agents
		OWNER TO arc;
`

var factsTable = `
	CREATE TABLE IF NOT EXISTS facts 
	(
		factID integer PRIMARY KEY NOT null,
		name varchar(255) NOT null,
		value varchar(255) NOT null,
		agent_id int, FOREIGN KEY (agent_id) 
		REFERENCES agents(uid), 
		CONSTRAINT uc_factID UNIQUE (factID, name)
	)
	WITH (
	 OIDS=FALSE
	);
	ALTER TABLE agents
		OWNER TO arc;
`

var Tables = [...]string{
	jobsTable,
	agentsTable,
	factsTable,
}

var InsertJobQuery = `INSERT INTO jobs(version,sender,requestid,"to",timeout,agent,action,payload,status) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) returning requestid;`
var GetAllJobsQuery = "SELECT * FROM jobs order by requestid"
var GetJobQuery = "SELECT * FROM jobs WHERE requestid=$1"
