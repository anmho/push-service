import * as dotenv from 'dotenv';
import * as path from 'path';
import { z } from 'zod';

dotenv.config({
  path: path.resolve(__dirname, '../lambda/.env.production'),
});

export const env = z
  .object({
    STRIPE_KEY: z.string(),
  })
  .parse(process.env);
