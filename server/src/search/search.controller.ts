import { Controller, Get, Query } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiQuery } from '@nestjs/swagger';
import { SearchService } from './search.service';

@ApiTags('Search')
@Controller('search')
export class SearchController {
  constructor(private readonly searchService: SearchService) {}

  @Get()
  @ApiOperation({ summary: 'Search for users or posts' })
  @ApiQuery({ name: 'query', required: true })
  @ApiQuery({ name: 'type', required: true, enum: ['users', 'posts'] })
  search(@Query('query') query: string, @Query('type') type: string) {
    // Basic validation handled by service or manual check, enum in Swagger is just doc unless using Pipes with Enum
    return this.searchService.search(query, type);
  }
}

