/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.bolingcavalry;

import org.apache.commons.lang3.StringUtils;
import org.apache.flink.api.common.functions.AggregateFunction;
import org.apache.flink.api.common.functions.MapFunction;
import org.apache.flink.api.java.functions.KeySelector;
import org.apache.flink.api.java.tuple.Tuple2;
import org.apache.flink.api.java.tuple.Tuple3;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.streaming.api.windowing.time.Time;
import org.apache.flink.streaming.connectors.wikiedits.WikipediaEditEvent;
import org.apache.flink.streaming.connectors.wikiedits.WikipediaEditsSource;

/**
 * Skeleton for a Flink Streaming Job.
 *
 * <p>For a tutorial how to write a Flink streaming application, check the
 * tutorials and examples on the <a href="http://flink.apache.org/docs/stable/">Flink Website</a>.
 *
 * <p>To package your application into a JAR file for execution, run
 * 'mvn clean package' on the command line.
 *
 * <p>If you change the name of the main class (with the public static void main(String[] args))
 * method, change the respective entry in the POM.xml file (simply search for 'mainClass').
 */
public class StreamingJob {

	public static void main(String[] args) throws Exception {
		// ????????????
		final StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

		env.addSource(new WikipediaEditsSource())
				//???????????????key??????
				.keyBy((KeySelector<WikipediaEditEvent, String>) wikipediaEditEvent -> wikipediaEditEvent.getUser())
				//???????????????5???
				.timeWindow(Time.seconds(15))
				//????????????????????????key????????????????????????
				.aggregate(new AggregateFunction<WikipediaEditEvent, Tuple3<String, Integer, StringBuilder>, Tuple3<String, Integer, StringBuilder>>() {
					@Override
					public Tuple3<String, Integer, StringBuilder> createAccumulator() {
						//??????ACC
						return new Tuple3<>("", 0, new StringBuilder());
					}

					@Override
					public Tuple3<String, Integer, StringBuilder> add(WikipediaEditEvent wikipediaEditEvent, Tuple3<String, Integer, StringBuilder> tuple3) {

						StringBuilder sbud = tuple3.f2;

						//????????????????????????????????????"Details ???"???????????????
						//?????????????????????????????????????????????????????????
						if(StringUtils.isBlank(sbud.toString())){
							sbud.append("Details : ");
						}else {
							sbud.append(" ");
						}

						//??????????????????????????????????????????
						return new Tuple3<>(wikipediaEditEvent.getUser(),
								wikipediaEditEvent.getByteDiff() + tuple3.f1,
								sbud.append(wikipediaEditEvent.getByteDiff()));
					}

					@Override
					public Tuple3<String, Integer, StringBuilder> getResult(Tuple3<String, Integer, StringBuilder> tuple3) {
						return tuple3;
					}

					@Override
					public Tuple3<String, Integer, StringBuilder> merge(Tuple3<String, Integer, StringBuilder> tuple3, Tuple3<String, Integer, StringBuilder> acc1) {
						//?????????????????????????????????
						return new Tuple3<>(tuple3.f0,
								tuple3.f1 + acc1.f1, tuple3.f2.append(acc1.f2));
					}
				})
				//???????????????????????????key????????????????????????????????????
				.map((MapFunction<Tuple3<String, Integer, StringBuilder>, String>) tuple3 -> tuple3.toString())
				//???????????????STDOUT
				.print();

		// ??????
		env.execute("Flink Streaming Java API Skeleton");
	}
}
