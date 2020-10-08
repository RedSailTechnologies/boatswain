import { TestBed } from '@angular/core/testing';

import { KrakenService } from './kraken.service';

describe('KrakenService', () => {
  let service: KrakenService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(KrakenService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
