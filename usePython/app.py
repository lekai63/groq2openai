from flask import Flask, request, jsonify
import requests
import os
import logging

app = Flask(__name__)

# 从环境变量获取代理服务器配置
default_groq_api_key= os.getenv('GROQ_API_KEY', 'your default groq api key')
proxy_url = os.getenv('PROXY_URL', 'http://127.0.0.1:7890')
proxies = {
    "http": proxy_url,
    "https": proxy_url,
}

# logging.basicConfig(level=logging.DEBUG)

@app.route('/v1/audio/transcriptions', methods=['POST'])
def transcribe_audio():
    try:
        # 打印请求信息
        # logging.debug('Received request headers: %s', request.headers)
        # logging.debug('Received request files: %s', request.files)
        # logging.debug('Received request form: %s', request.form)

        # 接收来自内网的请求
        openai_api_key = request.headers.get('Authorization').replace('Bearer ', '')  # 提取API Key
        file = request.files['file']
        model = request.form.get('model')

        # 设置默认值
        groq_api_key = openai_api_key if openai_api_key else default_groq_api_key
        groq_model = model if model else 'whisper-large-v3'

        groq_url = "https://api.groq.com/openai/v1/audio/transcriptions"

        files = {
            'file': (file.filename, file.stream, file.mimetype),
        }
        data = {
            'model': groq_model,
            'temperature': '0',
            'language': 'zh',
            'response_format': 'json',
            # response_format: Define the output response format.
                # Default is "json"
                # Set to "verbose_json" to receive timestamps for audio segments
                # Set to "text" to return a text response
                # formats vtt and srt are not supported
        }
        headers = {
            'Authorization': f'Bearer {groq_api_key}',
        }

        # 发送请求到Groq
        # logging.debug('Sending request to Groq with headers: %s, data: %s', headers, data)
        response = requests.post(groq_url, headers=headers, files=files, data=data, proxies=proxies)

        # 检查响应状态码
        response.raise_for_status()
        
        # 返回Groq的响应
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        logging.error('Request failed: %s', str(e))
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)