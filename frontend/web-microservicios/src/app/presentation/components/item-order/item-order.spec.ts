import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ItemOrder } from './item-order';

describe('ItemOrder', () => {
  let component: ItemOrder;
  let fixture: ComponentFixture<ItemOrder>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ItemOrder]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ItemOrder);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
