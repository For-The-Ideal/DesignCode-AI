
import { localUrl } from '~/config/index.js'
/**
 * 服务器获取数据 统一接口
 * url    请求接口
 * params 传参对象
 * isCookie 是否带上登录的cookie请求 默认是
 * method  请求方式  POST  GET  默认 POST
 */
export const getData = async (options = {}) => {
    let {
        url = '',
        req = {},
        params,
        isCookie = true,
        method,
    } = options;
    let reqJson = {};
    let result = {};
    const startTime = Date.now();
    try {
        let cookie = req.cookie || '';
        reqJson = {
            method,
            headers: {
                ...(isCookie && cookie ? { cookie } : {}),
            },
            baseURL: '', // 关键！
            timeout: 60e3, // 默认超时60秒
            credentials: 'include',
            mode: 'cors'
        };
        if (["get", "delete"].includes(method)) {
            reqJson.params = params;
        } else if (["post", "put"].includes(method)) {
            reqJson.body = params;
        }
        let baseURL = localUrl + url;

        result = await $fetch(baseURL, reqJson).then((res) => {
            return typeof res === 'string' ? JSON.parse(res) : res
        }).catch((error) => {
            return Promise.reject(error);
        });

        if(result.code !== 200){
            console.error('serverData 耗时：', Date.now() - startTime, '  地址：', url,'  参数：', JSON.stringify(reqJson), '返回数据：', JSON.stringify(result));
            return Promise.reject(result);
        }
        return result
    } catch (err) {
        const message = err ? err.message || err.data : err || '请求错误~！';
        console.error('serverData 耗时：', Date.now() - startTime, '  地址：', url,'  参数：', JSON.stringify(reqJson), '错误：', message);
        return {
            code: 400,
            error: 'serverData error',
            message,
            data: []
        };
    }
}
