import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewMyProducts } from './view-my-products';

describe('ViewMyProducts', () => {
  let component: ViewMyProducts;
  let fixture: ComponentFixture<ViewMyProducts>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewMyProducts]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewMyProducts);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
