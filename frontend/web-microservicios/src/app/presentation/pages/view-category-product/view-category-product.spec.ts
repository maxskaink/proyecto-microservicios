import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewCategoryProduct } from './view-category-product';

describe('ViewCategoryProduct', () => {
  let component: ViewCategoryProduct;
  let fixture: ComponentFixture<ViewCategoryProduct>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewCategoryProduct]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewCategoryProduct);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
