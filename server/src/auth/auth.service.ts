import { Injectable, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { InjectPinoLogger, PinoLogger } from 'nestjs-pino';
import * as bcrypt from 'bcryptjs';
import { PrismaService } from '../prisma.service';
import { LoginDto } from './dto/login.dto';
import { AuthResponseDto } from './dto/auth-response.dto';

@Injectable()
export class AuthService {
  constructor(
    private prisma: PrismaService,
    private jwtService: JwtService,
    @InjectPinoLogger(AuthService.name)
    private readonly logger: PinoLogger,
  ) {}

  async validateUserCredentials(username: string, pass: string): Promise<any> {
    const user = await this.prisma.user.findUnique({
      where: { username },
      select: {
        id: true,
        username: true,
        password: true,
        avatar: true,
        bio: true,
        following: {
          select: {
            followingId: true,
          },
        },
      },
    });

    if (!user) {
      this.logger.warn(`Login failed: user not found - ${username}`);
      throw new UnauthorizedException('Invalid credentials');
    }

    const isPasswordValid = await bcrypt.compare(pass, user.password);

    if (!isPasswordValid) {
      this.logger.warn(`Login failed: invalid password for user - ${username}`);
      throw new UnauthorizedException('Invalid credentials');
    }

    const { password, ...result } = user;
    return {
      ...result,
      following: user.following.map((f) => f.followingId),
    };
  }

  async login(loginDto: LoginDto): Promise<AuthResponseDto> {
    this.logger.info({ username: loginDto.username }, `Login attempt`);
    
    const user = await this.validateUserCredentials(
      loginDto.username,
      loginDto.password,
    );
    
    const payload = { userId: user.id };
    const token = this.jwtService.sign(payload);
    
    this.logger.info({ userId: user.id, username: user.username }, `Successful login`);
    
    return {
      token,
      user: {
        id: user.id,
        username: user.username,
        avatar: user.avatar || undefined,
        bio: user.bio || undefined,
        following: user.following,
      },
    };
  }
}

