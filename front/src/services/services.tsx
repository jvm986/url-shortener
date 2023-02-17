import axios, { AxiosResponse } from "axios";

const API_URL =
  process.env.NODE_ENV == "production"
    ? "https://l0xb0od05c.execute-api.eu-central-1.amazonaws.com/Prod/"
    : "http://127.0.0.1:5000/";

export async function getUrlInfo(url: string): Promise<AxiosResponse<any>> {
  const response = await axios.post(`${API_URL}/shorten`, {
    url: url,
  });
  return response;
}

export function validURL(str: string) {
  var pattern = new RegExp(
    "^(https?:\\/\\/)?" + // protocol
      "((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|" + // domain name
      "((\\d{1,3}\\.){3}\\d{1,3}))" + // OR ip (v4) address
      "(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*" + // port and path
      "(\\?[;&a-z\\d%_.~+=-]*)?" + // query string
      "(\\#[-a-z\\d_]*)?$",
    "i"
  ); // fragment locator
  return !!pattern.test(str);
}
