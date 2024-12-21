import React from "react";
import { List } from "antd";

const GameLog = ({ logs }) => (
  <List
    bordered
    dataSource={logs}
    renderItem={(log) => <List.Item>{log}</List.Item>}
    style={{ maxHeight: "200px", overflowY: "auto" }}
  />
);

export default GameLog;
