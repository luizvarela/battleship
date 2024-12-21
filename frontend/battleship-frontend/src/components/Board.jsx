import React from "react";
import "./Board.css";

const getCellColor = (value) => {
    switch (value) {
      case "ship":
        return "#666"; // Gray color for ships
      case "HIT":
        return "red";
      case "MISS":
        return "blue";
      default:
        return "white";
    }
  };

const Board = ({ data, onCellClick }) => {
  return (
    <div className="board" style={{ display: 'flex', flexDirection: 'column' }}>
      {data.map((row, rowIndex) => (
        <div key={rowIndex} className="board-row" style={{ display: 'flex', flexDirection: 'row' }}>
          {row.map((cell, colIndex) => (
            <div
              key={`${rowIndex}-${colIndex}`}
              className="board-cell"
              style={{ 
                backgroundColor: getCellColor(cell),
                border: '1px solid #999',
                width: '30px',
                height: '30px',
                cursor: 'pointer'
              }}
              onClick={() => onCellClick(rowIndex, colIndex)}
            />
          ))}
        </div>
      ))}
    </div>
  );
};

export default Board;
