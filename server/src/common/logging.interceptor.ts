import {
  Injectable,
  NestInterceptor,
  ExecutionContext,
  CallHandler,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';

@Injectable()
export class LoggingInterceptor implements NestInterceptor {
  constructor(
    @InjectPinoLogger(LoggingInterceptor.name)
    private readonly logger: PinoLogger,
  ) {
    this.logger.setContext('HTTP');
  }

  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    const request = context.switchToHttp().getRequest();
    const { method, url, body, query, params } = request;
    const userAgent = request.get('user-agent') || '';
    const ip = request.ip || request.connection.remoteAddress;
    const now = Date.now();

    // Логируем входящий запрос
    this.logger.info(
      {
        method,
        url,
        ip,
        userAgent,
        query: Object.keys(query).length > 0 ? query : undefined,
        params: Object.keys(params).length > 0 ? params : undefined,
        body: body && Object.keys(body).length > 0
          ? (body.password ? { ...body, password: '***' } : body)
          : undefined,
      },
      `Incoming request: ${method} ${url}`,
    );

    return next.handle().pipe(
      tap({
        next: (data) => {
          const response = context.switchToHttp().getResponse();
          const { statusCode } = response;
          const responseTime = Date.now() - now;
          this.logger.info(
            {
              method,
              url,
              statusCode,
              responseTime,
            },
            `Outgoing response: ${method} ${url} ${statusCode} - ${responseTime}ms`,
          );
        },
        error: (error) => {
          const responseTime = Date.now() - now;
          this.logger.error(
            {
              method,
              url,
              statusCode: error.status || 500,
              responseTime,
              error: error.message,
              stack: error.stack,
            },
            `Request failed: ${method} ${url} ${error.status || 500} - ${responseTime}ms`,
          );
        },
      }),
    );
  }
}

