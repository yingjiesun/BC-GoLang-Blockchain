package com.bcgeek.web.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import com.bcgeek.web.entity.BlockStruct;
import com.fasterxml.jackson.databind.ObjectMapper;

@Service
public class BlockChainService {

	@Autowired
	private RestTemplate restTemplate;

	public BlockStruct[] retrieveEntireBlocks() {

		try {
			String blocksInString = restTemplate.getForObject("http://localhost:8881", String.class);
			ObjectMapper mapper = new ObjectMapper();
			BlockStruct[] blockStructs = mapper.readValue(blocksInString, BlockStruct[].class);
			return blockStructs;
		} catch (Exception e) {
			e.printStackTrace();
		}

		return null;
	}
}
