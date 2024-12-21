// Ship component
import { motion } from "framer-motion";

const TILE_SIZE = 40;

const Ship = ({ coordinates, onPlace }) => {
  return (
    <motion.div
      className="ship"
      initial={{ scale: 0 }}
      animate={{ scale: 1 }}
      transition={{ duration: 0.5 }}
      style={{
        position: "absolute",
        top: coordinates.y * TILE_SIZE,
        left: coordinates.x * TILE_SIZE,
        width: TILE_SIZE,
        height: TILE_SIZE,
      }}
      onClick={() => onPlace(coordinates)}
    />
  );
};
