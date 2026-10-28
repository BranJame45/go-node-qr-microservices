import express from 'express';
import statisticsRoutes from './routes/statistics.js';
import cors from 'cors';

const app = express();

app.use(cors());
app.use(express.json());

app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

app.use('/api/statistics', statisticsRoutes);

const PORT = process.env.NODE_PORT || 3000;

if (process.env.NODE_ENV !== 'test') {
  app.listen(PORT, () => {
    console.log(`API Node corriendo en puerto ${PORT}`);
  });
}

export default app;
