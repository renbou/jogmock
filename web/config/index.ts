export function isProd(mode?: string) {
  if (mode === undefined) {
    return process.env.NODE_ENV === "production";
  }
  return mode === "production";
}

export function isDev(mode?: string) {
  return !isProd(mode);
}

export const duringDev = isDev();
export const duringProd = isProd();
