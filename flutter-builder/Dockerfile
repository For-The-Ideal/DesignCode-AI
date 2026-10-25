FROM ubuntu:22.04
ENV DEBIAN_FRONTEND=noninteractive

# 基础依赖 + Chrome(flutter web 构建/渲染需要) + git + curl
RUN apt-get update && apt-get install -y --no-install-recommends \
    git curl unzip xz-utils ca-certificates \
    wget gnupg2 \
    && rm -rf /var/lib/apt/lists/*

# 安装 coscli（COS 命令行 v5）
RUN curl -L -o /usr/local/bin/coscli \
    https://github.com/tencentyun/coscli/releases/latest/download/coscli-linux \
    && chmod +x /usr/local/bin/coscli

# 安装 Flutter stable（走国内镜像加速 pub）
ENV FLUTTER_HOME=/opt/flutter
RUN git clone -b stable https://github.com/flutter/flutter.git $FLUTTER_HOME \
    && echo 'export PATH=$PATH:$FLUTTER_HOME/bin:$FLUTTER_HOME/bin/cache/dart-sdk/bin' >> /etc/profile
ENV PATH="$PATH:/opt/flutter/bin:/opt/flutter/bin/cache/dart-sdk/bin"
ENV PUB_HOSTED_URL=https://pub.flutter-io.cn
ENV FLUTTER_STORAGE_BASE_URL=https://storage.flutter-io.cn

# 预热：创建模板工程并预编译一次(把 pub 依赖 & web 引擎缓存进镜像，省每次下载时间)
COPY template/ /opt/template/
RUN cd /opt/template && flutter config --no-analytics && flutter pub get || true

COPY build.sh /usr/local/bin/build.sh
RUN chmod +x /usr/local/bin/build.sh
ENTRYPOINT ["/usr/local/bin/build.sh"]