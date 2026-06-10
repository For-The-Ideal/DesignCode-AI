import { readBody, getRequestHeaders } from "h3";
import { getData } from "../utils/request.js";
import { getDecrypt } from "../utils/helps.js";

export default defineEventHandler(async (event) => {
  const body = await readBody(event);
  if (!body) {
    return { code: 400, message: "参数错误!", data: {} };
  }
  const req = getRequestHeaders(event);

  const {
    url = "",
    method = "post",
    params = {},
  } = req.client ? await readBody(event) : await getDecrypt(body.aes)

  return getData({
    method,
    url,
    params,
    req,
  });
});
