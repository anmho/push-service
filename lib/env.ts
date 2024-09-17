import * as dotenv from 'dotenv';
import * as path from 'path';
import { z } from 'zod';

dotenv.config({
  path: path.resolve(__dirname, '../lambda/.env.production'),
});

export const env = z
  .object({
  })
  .parse(process.env);
