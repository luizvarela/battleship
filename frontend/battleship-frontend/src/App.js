import React, { useState, useEffect } from "react";
import { Layout, message } from "antd";
import useWebSocket from "react-use-websocket";
import Board from "./components/Board";
import ShipPlacement from "./components/ShipPlacement";
import GameLog from "./components/GameLog";

const { Header, Content, Footer } = Layout;

const App = () => {
  const [playerBoard, setPlayerBoard] = useState(
    Array(10).fill(Array(10).fill("white"))
  );
  const [enemyBoard, setEnemyBoard] = useState(
    Array(10).fill(Array(10).fill("white"))
  );
  const [logs, setLogs] = useState([]);
  const [isConnected, setIsConnected] = useState(false);
  const [isAttackMode, setIsAttackMode] = useState(false);

  const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(
    "ws://localhost:8080/ws",
    {
      onOpen: () => {
        setIsConnected(true);
        addLog("Connected to the server.");
      },
      onClose: () => {
        setIsConnected(false);
        addLog("Disconnected from the server.");
      },
      onError: () => {
        message.error("WebSocket connection error.");
      },
      shouldReconnect: () => true, // Auto-reconnect on disconnection
    }
  );

  const [gameOver, setGameOver] = useState(false);

  useEffect(() => {
    if (lastJsonMessage?.type === "game_over") {
      setGameOver(true);
    }
  }, [lastJsonMessage]);

  useEffect(() => {
    if (lastJsonMessage) {
      handleServerMessage(lastJsonMessage);
    }
  }, [lastJsonMessage]);

  const addLog = (log) => setLogs((prev) => [...prev, log]);

  const handleServerMessage = (serverMsg) => {
    switch (serverMsg.type) {
      case "update_board":
        setPlayerBoard(serverMsg.board);
        break;
      case "error":
        message.error(serverMsg.message);
        addLog(`Error: ${serverMsg.message}`);
        break;
      case "enemy_update":
        setEnemyBoard(serverMsg.board);
        break;
      case "attack_result":
        const { board, result, row, col } = serverMsg;
        if (result === "HIT") {
          setEnemyBoard((prev) =>
            prev.map((r, rowIndex) =>
              rowIndex === row
                ? r.map((c, colIndex) => (colIndex === col ? "HIT" : c))
                : r
            )
          );
        } else if (result === "MISS") {
          setEnemyBoard((prev) =>
            prev.map((r, rowIndex) =>
              rowIndex === row
                ? r.map((c, colIndex) => (colIndex === col ? "MISS" : c))
                : r
            )
          );
        }
        break;
      case "log":
        addLog(serverMsg.text);
        break;
      case "game_over":
        addLog(`Game Over! ${serverMsg.winner} wins.`);
        break;
      default:
        console.warn("Unknown message type:", serverMsg);
    }
  };

  const handlePlaceShip = (ship, orientation, row, col) => {
    if (!isConnected) {
      message.warning("Not connected to the server.");
      return;
    }
    sendJsonMessage({ type: "place_ship", ship, orientation, row, col });
    addLog(`Placing ${ship} (${orientation}) at position (${row}, ${col}).`);
  };

  const handleAttack = (row, col) => {
    if (!isConnected) {
      message.warning("Not connected to the server.");
      return;
    }

    if (enemyBoard[row][col] === "hit" || enemyBoard[row][col] === "miss") {
      message.warning("This cell has already been attacked!");
      return;
    }

    sendJsonMessage({ type: "attack", row, col });
    addLog(`Attacking (${row}, ${col}).`);
  };

  const handleCellClick = (row, col) => {
    if (isAttackMode && !gameOver) {
      handleAttack(row, col);
    }
  };

  return (
    <Layout style={{ height: "100vh" }}>
      <Header style={{ color: "white" }}>
        Battleship Game {isConnected ? "(Connected)" : "(Disconnected)"}
      </Header>
      <Content style={{ padding: "20px" }}>
        <ShipPlacement onPlaceShip={handlePlaceShip} />
        <div style={{ marginBottom: "20px" }}>
          <button
            onClick={() => setIsAttackMode(!isAttackMode)}
            style={{
              padding: "8px 16px",
              backgroundColor: isAttackMode ? "#ff4d4f" : "#1890ff",
              color: "white",
              border: "none",
              borderRadius: "4px",
              cursor: "pointer",
            }}
          >
            {isAttackMode ? "Cancel Attack" : "Start Attack"}
          </button>
        </div>
        <div style={{ display: "flex", justifyContent: "space-between" }}>
          <div>
            <h3>Your Board</h3>
            <Board data={playerBoard} onCellClick={() => {}} />
          </div>
          <div>
            <h3>Enemy Board {isAttackMode && "(Click to Attack)"}</h3>
            <Board data={enemyBoard} onCellClick={handleCellClick} />
          </div>
        </div>

        {gameOver && (
          <div className="game-over">
            <h1>Game Over!</h1>
            <p>{`Winner: ${lastJsonMessage?.winner}`}</p>
          </div>
        )}

        <GameLog logs={logs} />
      </Content>
      <Footer style={{ textAlign: "center" }}>Battleship Game ©2024</Footer>
    </Layout>
  );
};

export default App;
