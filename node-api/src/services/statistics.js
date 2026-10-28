// Servicio de cálculo de estadísticas
export function calculateStatistics(matrices) {
  const result = {
    global: {
      max: 0,
      min: 0,
      average: 0,
      sum: 0,
      isDiagonal: false
    },
    perMatrix: {}
  };

  if (!matrices || typeof matrices !== 'object') {
    return result;
  }

  const allValues = [];
  let anyDiagonal = false;

  for (const [key, matrix] of Object.entries(matrices)) {
    if (!Array.isArray(matrix)) {
      result.perMatrix[key] = {
        max: 0,
        min: 0,
        average: 0,
        sum: 0,
        isDiagonal: false
      };
      continue;
    }

    const values = [];
    for (let i = 0; i < matrix.length; i++) {
      const row = matrix[i];
      if (!Array.isArray(row)) continue;
      for (let j = 0; j < row.length; j++) {
        const val = Number(row[j]);
        if (!isNaN(val)) {
          values.push(val);
          allValues.push(val);
        }
      }
    }

    if (values.length === 0) {
      result.perMatrix[key] = {
        max: 0,
        min: 0,
        average: 0,
        sum: 0,
        isDiagonal: false
      };
      continue;
    }

    let sum = 0;
    let max = values[0];
    let min = values[0];
    for (let i = 0; i < values.length; i++) {
      const v = values[i];
      sum += v;
      if (v > max) max = v;
      if (v < min) min = v;
    }
    const average = sum / values.length;

    // Verificar si es diagonal: SOLO si es cuadrada y todo fuera de diagonal |valor| < 1e-9
    let isDiagonal = false;
    if (matrix.length > 0 && Array.isArray(matrix[0]) && matrix.length === matrix[0].length) {
      isDiagonal = true;
      for (let i = 0; i < matrix.length; i++) {
        const row = matrix[i];
        if (!Array.isArray(row)) {
          isDiagonal = false;
          break;
        }
        if (row.length !== matrix.length) {
          isDiagonal = false;
          break;
        }
        for (let j = 0; j < row.length; j++) {
          if (i !== j) {
            const val = Number(row[j]);
            if (Math.abs(val) >= 1e-9) {
              isDiagonal = false;
              break;
            }
          }
        }
        if (!isDiagonal) break;
      }
    }

    if (isDiagonal) {
      anyDiagonal = true;
    }

    result.perMatrix[key] = {
      max,
      min,
      average,
      sum,
      isDiagonal
    };
  }

  if (allValues.length === 0) {
    result.global = {
      max: 0,
      min: 0,
      average: 0,
      sum: 0,
      isDiagonal: anyDiagonal
    };
    return result;
  }

  let globalSum = 0;
  let globalMax = allValues[0];
  let globalMin = allValues[0];
  for (let i = 0; i < allValues.length; i++) {
    const v = allValues[i];
    globalSum += v;
    if (v > globalMax) globalMax = v;
    if (v < globalMin) globalMin = v;
  }
  result.global = {
    max: globalMax,
    min: globalMin,
    average: globalSum / allValues.length,
    sum: globalSum,
    isDiagonal: anyDiagonal
  };

  return result;
}
