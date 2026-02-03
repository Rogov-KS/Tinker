import { ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty, IsString } from 'class-validator';

export class CreatePostDto {
  @ApiProperty({ example: 'Hello world!' })
  @IsString()
  @IsNotEmpty()
  content!: string;
}

