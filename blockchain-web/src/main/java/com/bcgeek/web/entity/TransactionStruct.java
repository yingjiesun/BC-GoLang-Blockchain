package com.bcgeek.web.entity;

import com.fasterxml.jackson.annotation.JsonProperty;

public class TransactionStruct {

	@JsonProperty("TransactionId")
	private String transactionId;
	
	@JsonProperty("Timestamp")
	private String timestamp;
	
	public String getTransactionId() {
		return transactionId;
	}
	public void setTransactionId(String transactionId) {
		this.transactionId = transactionId;
	}
	public String getTimestamp() {
		return timestamp;
	}
	public void setTimestamp(String timestamp) {
		this.timestamp = timestamp;
	}
	
}
