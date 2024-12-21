import React, { useState } from "react";
import { Button, InputNumber } from "antd";

const AttackPanel = ({ onAttack }) => {
  const [row, setRow] = useState(0);
  const [col, setCol] = useState(0);

  const handleAttack = () => {
    onAttack(row, col);
  };

  return (
    <div style={{ marginBottom: "20px" }}>
      <InputNumber
        min={0}
        max={9}
        value={row}
        onChange={setRow}
        style={{ marginRight: "10px" }}
        placeholder="Row"
      />
      <InputNumber
        min={0}
        max={9}
        value={col}
        onChange={setCol}
        style={{ marginRight: "10px" }}
        placeholder="Column"
      />
      <Button type="primary" onClick={handleAttack}>
        Attack
      </Button>
    </div>
  );
};

export default AttackPanel;
