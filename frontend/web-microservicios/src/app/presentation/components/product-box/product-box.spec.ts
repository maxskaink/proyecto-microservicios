import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProductBox } from './product-box';

describe('ProductBox', () => {
  let component: ProductBox;
  let fixture: ComponentFixture<ProductBox>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ProductBox]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ProductBox);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
