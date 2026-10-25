# 换成国内 DaoCloud 镜像源，秒拉取
FROM docker.m.daocloud.io/ubuntu:22.04
ENV DEBIAN_FRONTEND=noninteractive

# 基础依赖 + Chrome + git + curl
RUN apt-get update && apt-get install -y --no-install-recommends \
    git curl unzip xz-utils ca-certificates \
    wget gnupg2 \
    && rm -rf /var/lib/apt/lists/*

# 安装 coscli（使用腾讯云官方国内镜像地址）
RUN curl -L -o /usr/local/bin/coscli \
    https://cosbrowser.cloud.tencent.com/software/coscli/coscli-linux-amd64 \
    && chmod +x /usr/local/bin/coscli

# 安装 Flutter stable（**关键修改**：换成了国内的 Gitee 镜像仓库，飞快！）
ENV FLUTTER_HOME=/opt/flutter
RUN git clone -b stable https://gitee.com/mirrors/Flutter.git $FLUTTER_HOME \
    && echo 'export PATH=$PATH:$FLUTTER_HOME/bin:$FLUTTER_HOME/bin/cache/dart-sdk/bin' >> /etc/profile
ENV PATH="$PATH:/opt/flutter/bin:/opt/flutter/bin/cache/dart-sdk/bin"

# 走国内 Flutter 镜像
ENV PUB_HOSTED_URL=https://pub.flutter-io.cn
ENV FLUTTER_STORAGE_BASE_URL=https://storage.flutter-io.cn

# 预热
COPY template/ /opt/template/
RUN cd /opt/template && flutter config --no-analytics && flutter pub get || true

COPY build.sh /usr/local/bin/build.sh
RUN chmod +x /usr/local/bin/build.sh
ENTRYPOINT ["/usr/local/bin/build.sh"]