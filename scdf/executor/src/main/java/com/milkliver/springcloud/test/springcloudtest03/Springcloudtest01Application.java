package com.milkliver.springcloud.test.springcloudtest03;

import org.springframework.boot.ExitCodeGenerator;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Springcloudtest01Application {
//	public class Springcloudtest01Application implements ExitCodeGenerator {

	public static void main(String[] args) {
		SpringApplication.run(Springcloudtest01Application.class, args);
	}

//	@Override
//	public int getExitCode() {
//		// TODO Auto-generated method stub
//		return 0;
//	}

}
