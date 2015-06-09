package db

var jobsTable = `
	CREATE TABLE IF NOT EXISTS jobs
	(
		version integer,
		sender text,
		requestid text PRIMARY KEY NOT NULL,
		"to" text,
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
