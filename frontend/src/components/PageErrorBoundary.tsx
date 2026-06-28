import { Alert, Button } from "antd";
import { Component, type ErrorInfo, type ReactNode } from "react";

type Props = {
  children: ReactNode;
};

type State = {
  error?: Error;
};

export default class PageErrorBoundary extends Component<Props, State> {
  state: State = {};

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error("Page render failed", error, info);
  }

  render() {
    if (this.state.error) {
      return (
        <Alert
          type="error"
          showIcon
          message="页面渲染失败"
          description={this.state.error.message || "当前页面数据格式异常，请刷新后重试。"}
          action={<Button onClick={() => this.setState({ error: undefined })}>重试</Button>}
        />
      );
    }
    return this.props.children;
  }
}
