import request from 'supertest';
import app from '../src/index.js';
import jwt from 'jsonwebtoken';

const JWT_SECRET = process.env.JWT_SECRET || 'change-me-super-secret';

function generateToken() {
  return jwt.sign({ user: 'test' }, JWT_SECRET, { expiresIn: '1h' });
}

describe('Statistics API', () => {
  test('POST /api/statistics without token returns 401', async () => {
    const res = await request(app)
      .post('/api/statistics')
      .send({ matrices: { Q: [[1, 2]], R: [[3, 4]] } });
    expect(res.statusCode).toBe(401);
  });

  test('POST /api/statistics with valid token returns statistics', async () => {
    const token = generateToken();
    const res = await request(app)
      .post('/api/statistics')
      .set('Authorization', `Bearer ${token}`)
      .send({ matrices: { Q: [[1, 2]], R: [[3, 4]] } });
    expect(res.statusCode).toBe(200);
    expect(res.body).toHaveProperty('global');
    expect(res.body).toHaveProperty('perMatrix');
    expect(res.body.perMatrix).toHaveProperty('Q');
    expect(res.body.perMatrix).toHaveProperty('R');
  });

  test('Calculates max, min, sum, average correctly', async () => {
    const token = generateToken();
    const res = await request(app)
      .post('/api/statistics')
      .set('Authorization', `Bearer ${token}`)
      .send({ matrices: { Q: [[1, 2]], R: [[3]] } });
    expect(res.statusCode).toBe(200);
    expect(res.body.global.sum).toBe(6);
    expect(res.body.global.max).toBe(3);
    expect(res.body.global.min).toBe(1);
    expect(res.body.global.average).toBe(2);
  });

  test('isDiagonal detects diagonal matrix', async () => {
    const token = generateToken();
    const res = await request(app)
      .post('/api/statistics')
      .set('Authorization', `Bearer ${token}`)
      .send({ matrices: { Q: [[1, 0], [0, 1]], R: [[2, 1], [0, 3]] } });
    expect(res.statusCode).toBe(200);
    expect(res.body.perMatrix.Q.isDiagonal).toBe(true);
    expect(res.body.global.isDiagonal).toBe(true);
  });

  test('isDiagonal false for non-diagonal', async () => {
    const token = generateToken();
    const res = await request(app)
      .post('/api/statistics')
      .set('Authorization', `Bearer ${token}`)
      .send({ matrices: { Q: [[1, 2], [0, 1]], R: [[1, 0], [0, 1]] } });
    expect(res.statusCode).toBe(200);
    expect(res.body.perMatrix.Q.isDiagonal).toBe(false);
  });

  test('isDiagonal false for rectangular matrix', async () => {
    const token = generateToken();
    const res = await request(app)
      .post('/api/statistics')
      .set('Authorization', `Bearer ${token}`)
      .send({ matrices: { Q: [[5, 0], [0, 3], [0, 0]], R: [[1, 0], [0, 1]] } });
    expect(res.statusCode).toBe(200);
    expect(res.body.perMatrix.Q.isDiagonal).toBe(false);
  });
});
