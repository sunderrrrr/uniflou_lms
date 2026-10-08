export class ApiError extends Error {
  code: string;

  constructor(code: string, message: string) {
    super(message);
    this.code = code;
  }
}

export async function parseApiError(res: Response): Promise<never> {
  try {
    const data = (await res.json()) as { error?: { code?: string; message?: string } };
    throw new ApiError(data.error?.code ?? 'UNKNOWN', data.error?.message ?? 'Ошибка запроса');
  } catch (e) {
    if (e instanceof ApiError) throw e;
    throw new ApiError('UNKNOWN', res.statusText || 'Ошибка запроса');
  }
}
