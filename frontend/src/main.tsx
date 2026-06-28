import React from "react";
import ReactDOM from "react-dom/client";
import { ConfigProvider, theme } from "antd";
import zhCN from "antd/locale/zh_CN";
import App from "./App";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ConfigProvider
      locale={zhCN}
      theme={{
        algorithm: theme.darkAlgorithm,
        token: {
          colorPrimary: "#0A84FF",
          borderRadius: 8,
          colorBgBase: "#08090B",
          colorBgContainer: "rgba(255, 255, 255, 0.08)",
          colorBgElevated: "rgba(20, 22, 28, 0.86)",
          colorBorder: "rgba(255, 255, 255, 0.16)",
          colorBorderSecondary: "rgba(255, 255, 255, 0.09)",
          colorText: "rgba(255, 255, 255, 0.94)",
          colorTextSecondary: "rgba(235, 238, 245, 0.68)",
          colorTextTertiary: "rgba(235, 238, 245, 0.48)",
          fontFamily:
            '-apple-system, BlinkMacSystemFont, "SF Pro Display", "SF Pro Text", "Segoe UI", sans-serif',
          controlHeight: 40,
          controlHeightLG: 46,
          controlHeightSM: 32,
          boxShadow:
            "0 18px 48px rgba(0, 0, 0, 0.32), inset 0 1px 0 rgba(255, 255, 255, 0.14)",
        },
        components: {
          Card: {
            headerFontSize: 16,
            paddingLG: 22,
            borderRadiusLG: 8,
          },
          Table: {
            headerBg: "rgba(255, 255, 255, 0.08)",
            rowHoverBg: "rgba(255, 255, 255, 0.055)",
          },
          Modal: {
            contentBg: "rgba(18, 20, 24, 0.78)",
          },
          Input: {
            colorBgContainer: "rgba(255, 255, 255, 0.07)",
          },
          Select: {
            colorBgContainer: "rgba(255, 255, 255, 0.07)",
          },
          Button: {
            borderRadius: 999,
            primaryShadow: "0 10px 24px rgba(0, 113, 227, 0.32)",
          },
        },
      }}
    >
      <App />
    </ConfigProvider>
  </React.StrictMode>
);
