CREATE TABLE IF NOT EXISTS request_logs (
 	id BIGINT PRIMARY KEY AUTO_INCREMENT,
 	user_id VARCHAR(64),
	api_key VARCHAR(128),
	model VARCHAR(64),
	prompt TEXT,
	response TEXT,
	input_tokens INT DEFAULT 0,
	output_tokens INT DEFAULT 0,
	latency_ms INT,
	status VARCHAR(32),
 	error_message TEXT,
 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
