package com.bcgeek.web.config.security;

import org.springframework.context.annotation.Configuration;
import org.springframework.core.annotation.Order;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;

@Configuration
@Order(1)
public class CommonWebSecurityConfiguration extends WebSecurityConfigurerAdapter {

	// @formatter:off
	private final String[] WHITE_LIST = new String[] {
			"/css/**", 
			"/webjars/**", 
			"/img/**", 
			"/" };

	@Override
	protected void configure(HttpSecurity http) throws Exception {
		http
			.authorizeRequests()
				.antMatchers(WHITE_LIST)
					.permitAll()
					.anyRequest().authenticated();
	}
	// @formatter:on
}