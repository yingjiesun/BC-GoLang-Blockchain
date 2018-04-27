package com.bcgeek.web.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.RequestMapping;

import com.bcgeek.web.entity.BlockStruct;
import com.bcgeek.web.service.BlockChainService;

@Controller
public class HomeController {

	@Autowired
	BlockChainService blockchainService;

	@RequestMapping("/")
	public String home(Model model) {
		BlockStruct[] blocks = blockchainService.retrieveEntireBlocks();
		model.addAttribute("blocks", blocks);
		return "index";
	}
}
