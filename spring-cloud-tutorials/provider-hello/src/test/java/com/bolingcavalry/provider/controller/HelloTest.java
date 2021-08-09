package com.bolingcavalry.provider.controller;

import lombok.extern.slf4j.Slf4j;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;

import static org.hamcrest.Matchers.containsString;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.content;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

/**
 * @Description: 单元测试类
 * @author: willzhao E-mail: zq2599@gmail.com
 * @date: 2020/8/9 23:55
 */
@SpringBootTest
@AutoConfigureMockMvc
@Slf4j
class HelloTest {

    @Autowired
    private MockMvc mvc;

    @Test
    void hello() throws Exception {
        String responseString = mvc.perform(MockMvcRequestBuilders.get("/hello/str").accept(MediaType.APPLICATION_JSON))
                .andExpect(status().isOk())
                .andExpect(content().string(containsString(Hello.HELLO_PREFIX)))
                .andDo(print())
                .andReturn()
                .getResponse()
                .getContentAsString();

        log.info("response in junit test :\n" + responseString);
    }
}