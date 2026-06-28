import axios from "axios";
import { API_BASE } from "../config";

const client = axios.create({
  baseURL: API_BASE,
  timeout: 15000,
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

client.interceptors.response.use(
  (response) => response,
  (error) => {
    const apiMessage = error?.response?.data?.message;
    const fallback = error?.message || "请求失败";
    return Promise.reject(new Error(apiMessage || fallback));
  }
);

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export default client;
