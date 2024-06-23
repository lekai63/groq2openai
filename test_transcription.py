import requests

# 配置测试音频文件和服务器地址
file_path = 'sample.wav'
server_url = 'http://127.0.0.1:8000/v1/audio/transcriptions'

# 配置OpenAI API密钥和模型
openai_api_key = ''
model = 'whisper-large-v3'

# 打开音频文件
with open(file_path, 'rb') as f:
    files = {
        'file': f,
    }
    headers = {
        'Authorization': openai_api_key,
    }
    data = {
        'model': model,
    }

    # 发送POST请求到Flask应用
    response = requests.post(server_url, headers=headers, files=files, data=data)

# 打印响应
print(response.status_code)
try:
    print(response.json())
except requests.exceptions.JSONDecodeError:
    print("Response content is not in JSON format.")
    print(response.text)