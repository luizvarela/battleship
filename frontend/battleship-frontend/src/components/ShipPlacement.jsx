import React, { useState } from "react";
import { Select, Button, InputNumber, Space } from "antd";

const ShipPlacement = ({ onPlaceShip }) => {
  const [selectedShip, setSelectedShip] = useState("carrier");
  const [orientation, setOrientation] = useState("horizontal");
  const [row, setRow] = useState(0);
  const [col, setCol] = useState(0);

  const ships = [
    { value: "carrier", label: "Carrier (5)" },
    { value: "battleship", label: "Battleship (4)" },
    { value: "cruiser", label: "Cruiser (3)" },
    { value: "submarine", label: "Submarine (3)" },
    { value: "destroyer", label: "Destroyer (2)" },
  ];

  const handlePlacement = () => {
    onPlaceShip(selectedShip, orientation, row, col);
  };

  return (
    <div style={{ marginBottom: 20 }}>
      <Space>
        <Select
          value={selectedShip}
          onChange={setSelectedShip}
          style={{ width: 150 }}
        >
          {ships.map((ship) => (
            <Select.Option key={ship.value} value={ship.value}>
              {ship.label}
            </Select.Option>
          ))}
        </Select>

        <Select
          value={orientation}
          onChange={setOrientation}
          style={{ width: 120 }}
        >
          <Select.Option value="horizontal">Horizontal</Select.Option>
          <Select.Option value="vertical">Vertical</Select.Option>
        </Select>

        <InputNumber
          min={0}
          max={9}
          value={row}
          onChange={setRow}
          placeholder="Row (0-9)"
        />
        <InputNumber
          min={0}
          max={9}
          value={col}
          onChange={setCol}
          placeholder="Col (0-9)"
        />

        <Button type="primary" onClick={handlePlacement}>
          Place Ship
        </Button>
      </Space>
    </div>
  );
};

export default ShipPlacement;
