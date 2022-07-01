package com.milkliver.springcloud.test.springcloudtest03.Tasks;

import java.io.BufferedInputStream;
import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.URL;
import java.text.SimpleDateFormat;
import java.util.Base64;
import java.util.Date;
import java.util.HashMap;
import java.util.Locale;
import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.cloud.task.configuration.EnableTask;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import com.milkliver.springcloud.test.springcloudtest03.utils.SendRequest;

@Configuration
@EnableTask
public class Task01 {

	private static final Logger log = LoggerFactory.getLogger(Task01.class);

	final Base64.Decoder decoder = Base64.getDecoder();
	final Base64.Encoder encoder = Base64.getEncoder();

	@Value("${scdf.server.runtask.api.request.url:#{null}}")
	String scdfServerRuntaskApiRequestUrl;

	@Value("${scdf.server.runtask.api.request.hostname:#{null}}")
	String scdfServerRuntaskApiRequestHostname;

	@Value("${scdf.server.runtask.api.request.method:GET}")
	String scdfServerRuntaskApiRequestMethod;

	@Value("${scdf.server.runtask.api.request.connect-time-out:2000}")
	int scdfServerRuntaskApiRequestConnectTimeOut;

	@Value("${scdf.server.runtask.api.request.read-time-out:2000}")
	int scdfServerRuntaskApiRequestReadTimeOut;

	@Value("${scdf.server.runtask.api.request.enable-https:false}")
	Boolean scdfServerRuntaskApiRequestEnableHttps;

	@Value("${system.command}")
	String systemCommandBase64;

	@Autowired
	SendRequest sendRequest;

	@Value("${spring.cloud.task.executionid:#{null}}")
	String taskid;

	@Bean
	public CommandLineRunner commandLineRunner() {
		return args -> {

			log.info("SCDF executor CommandLineRunner ...");
			if (taskid != null) {
				log.info("taskid: " + String.valueOf(taskid) + " is running ...");
			} else {
				log.info("taskid: null is running ...");
			}

			try {

				log.info("CommandBase64: " + systemCommandBase64);
				String systemCommand = new String(decoder.decode(systemCommandBase64));
				log.info("Command: " + systemCommand);
				log.info(
						"==================================================start==================================================");

				Process process = Runtime.getRuntime().exec(systemCommand);
//				Process process = Runtime.getRuntime().exec("powershell ls");
//				Process process = Runtime.getRuntime().exec("notepad.exe");
//				Process process = Runtime.getRuntime().exec(
//						"java -jar D:\\JavaProjects\\TestProject03\\springcloud-test-task01\\externalProgramFiles\\java-job01.jar");
//				Process process = Runtime.getRuntime().exec("python D:\\JavaProjects\\TestProject03\\springcloud-test-task01\\externalProgramFiles\\test.py");

//				============================================================================================
				StringBuilder execCmdRes = new StringBuilder();

				BufferedReader bufferedReader = new BufferedReader(new InputStreamReader(process.getInputStream()));
				String line;
				while ((line = bufferedReader.readLine()) != null) {
					execCmdRes.append(line);
					execCmdRes.append("\r\n");
				}
//				============================================================================================

				log.info(execCmdRes.toString());

//				log.info("exitVaule: " + String.valueOf(process.exitValue()));
				log.info("waitFor: " + String.valueOf(process.waitFor()));

				Map<String, Object> jsonMap = new HashMap<String, Object>();

				jsonMap.put("message", taskid);
//				jsonMap.put("message", "1084");

				if (process.waitFor() != 0) {
					log.info("task is failed");
				} else {
					Map<String, Object> sendSuccessMsgRes = new HashMap<String, Object>();
					if (scdfServerRuntaskApiRequestEnableHttps) {
						sendSuccessMsgRes = sendRequest.https(scdfServerRuntaskApiRequestUrl,
								scdfServerRuntaskApiRequestHostname, scdfServerRuntaskApiRequestMethod,
								scdfServerRuntaskApiRequestConnectTimeOut, scdfServerRuntaskApiRequestReadTimeOut,
								jsonMap);
					} else {
						sendSuccessMsgRes = sendRequest.http(scdfServerRuntaskApiRequestUrl,
								scdfServerRuntaskApiRequestHostname, scdfServerRuntaskApiRequestMethod,
								scdfServerRuntaskApiRequestConnectTimeOut, scdfServerRuntaskApiRequestReadTimeOut,
								jsonMap);
					}

					log.info("response content: " + String.valueOf(sendSuccessMsgRes.get("responseContent")));
					log.info("task is success");
				}

				process.destroy();
				log.info(
						"===================================================end===================================================");

			} catch (Exception e) {
				log.error(e.getMessage());
				for (StackTraceElement elem : e.getStackTrace()) {
					log.error(elem.toString());
				}
			}
			log.info("CommandLineRunner finish");
			if (taskid != null) {
				log.info("taskid: " + String.valueOf(taskid) + " is finished ...");
			} else {
				log.info("taskid: null is finished ...");
			}
			log.info("SCDF executor CommandLineRunner finish");
		};
	}
}
