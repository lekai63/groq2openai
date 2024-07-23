FROM python:3.12-slim AS build

# 设置 pip 源为中科大源，并升级 pip
RUN pip config set global.index-url https://mirrors.ustc.edu.cn/pypi/web/simple && \
    pip install --upgrade pip

# 复制 requirements.txt 并安装依赖
COPY ./requirements.txt /requirements.txt
RUN pip install --timeout 30 --user --no-cache-dir --no-warn-script-location -r /requirements.txt

FROM python:3.12-slim

# 设置环境变量
ENV LOCAL_PKG="/root/.local"

# 复制已安装的依赖和时区信息
COPY --from=build ${LOCAL_PKG} ${LOCAL_PKG}
COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 设置符号链接并配置时区
RUN ln -sf ${LOCAL_PKG}/bin/* /usr/local/bin/ && echo "Asia/Shanghai" > /etc/timezone

# 设置工作目录
WORKDIR /app

# 复制应用代码到工作目录
COPY . /app

# 暴露应用运行的端口
EXPOSE 8000

# 运行 Flask 应用
# CMD ["python", "app.py"]
CMD ["gunicorn", "-w", "4", "-b", "0.0.0.0:8000", "app:app"]
