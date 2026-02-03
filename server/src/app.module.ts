import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { LoggerModule } from 'nestjs-pino';
import { AuthModule } from './auth/auth.module';
import { UsersModule } from './users/users.module';
import { PostsModule } from './posts/posts.module';
import { CommentsModule } from './comments/comments.module';
import { FeedModule } from './feed/feed.module';
import { SearchModule } from './search/search.module';
import { PrismaService } from './prisma.service';
import { LoggingInterceptor } from './common/logging.interceptor';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    LoggerModule.forRoot({
      pinoHttp: {
        // Отключаем автоматическое логирование HTTP запросов, используем наш LoggingInterceptor
        autoLogging: false,
        useLevel: 'silent', // Полностью отключаем HTTP логирование
        transport:
          process.env.NODE_ENV !== 'production'
            ? {
                target: 'pino-pretty',
                options: {
                  colorize: true,
                  singleLine: false,
                  translateTime: 'HH:MM:ss Z',
                  ignore: 'pid,hostname',
                },
              }
            : undefined,
        level: process.env.LOG_LEVEL || 'info',
      },
    }),
    AuthModule,
    UsersModule,
    PostsModule,
    CommentsModule,
    FeedModule,
    SearchModule,
  ],
  controllers: [],
  providers: [PrismaService, LoggingInterceptor],
})
export class AppModule {}
