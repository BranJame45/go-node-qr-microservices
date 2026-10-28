import express from 'express';
import { calculateStatistics } from '../services/statistics.js';
import { authMiddleware } from '../middleware/auth.js';

const router = express.Router();

// POST /api/statistics - requiere JWT
router.post('/', authMiddleware, (req, res) => {
  try {
    const { matrices } = req.body;
    if (!matrices || typeof matrices !== 'object') {
      return res.status(400).json({ error: 'Body inválido' });
    }
    const stats = calculateStatistics(matrices);
    res.json(stats);
  } catch (err) {
    res.status(400).json({ error: 'Body inválido' });
  }
});

export default router;
