import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HistoryOrdersUser } from './history-orders-user';

describe('HistoryOrdersUser', () => {
  let component: HistoryOrdersUser;
  let fixture: ComponentFixture<HistoryOrdersUser>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HistoryOrdersUser]
    })
    .compileComponents();

    fixture = TestBed.createComponent(HistoryOrdersUser);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
