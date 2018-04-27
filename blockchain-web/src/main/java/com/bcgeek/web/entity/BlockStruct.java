package com.bcgeek.web.entity;

import com.fasterxml.jackson.annotation.JsonProperty;

public class BlockStruct {

	@JsonProperty("Hash")
	private String hash;
	
	@JsonProperty("Index")
	private Long index;
	
	@JsonProperty("Nounce")
	private Long nounce;
	
	@JsonProperty("PrevHash")
	private String prevHash;
	
	@JsonProperty("Timestamp")
	private String timestamp;
	
	@JsonProperty("Transactions")
	private TransactionStruct[] transactionStructs;
	
	public String getHash() {
		return hash;
	}
	public void setHash(String hash) {
		this.hash = hash;
	}
	public Long getIndex() {
		return index;
	}
	public void setIndex(Long index) {
		this.index = index;
	}
	public Long getNounce() {
		return nounce;
	}
	public void setNounce(Long nounce) {
		this.nounce = nounce;
	}
	public String getPrevHash() {
		return prevHash;
	}
	public void setPrevHash(String prevHash) {
		this.prevHash = prevHash;
	}
	public String getTimestamp() {
		return timestamp;
	}
	public void setTimestamp(String timestamp) {
		this.timestamp = timestamp;
	}
	public TransactionStruct[] getTransactionStructs() {
		return transactionStructs;
	}
	public void setTransactionStructs(TransactionStruct[] transactionStructs) {
		this.transactionStructs = transactionStructs;
	}
	
	
}
