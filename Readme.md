
## 背景

国内很多语音转文字app或陪伴机器人默认使用whisper-medium或whisper-small，但对中文的识别率还是有所欠缺。
而whisper-large-v3在groq上运行速度非常快，且对中文的识别率也好的多，所以将groq的whisper-large-v3转换为openai格式，方便其他程序以openai_api格式调用.

## 调用方法

以火火兔folotoy配置使用本项目为例：

```yml
folotoy: #火火兔的配置
  environment:
      STT_TYPE: openai-whisper
      OPENAI_WHISPER_KEY: gsk_xxxx #填入groq key
      OPENAI_WHISPER_MODEL: whisper-large-v3
      OPENAI_WHISPER_API_BASE: http://groq:8000/v1

  groq-proxy: #本项目docker
    image: lekai/groq2openai:whisper #请自行编译
    container_name: groq
    ports:
      - "8009:8000"
    environment:
      # - HOST=127.0.0.1
      PROXY_URL: "http://mihomo:7890" # 通过代理连接groq
    restart: unless-stopped
```

本项目全程由chatgpt创作，感谢chatgpt。
