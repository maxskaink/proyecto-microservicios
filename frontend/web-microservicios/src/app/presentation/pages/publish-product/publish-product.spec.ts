import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PublishProduct } from './publish-product';

describe('PublishProduct', () => {
  let component: PublishProduct;
  let fixture: ComponentFixture<PublishProduct>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PublishProduct]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PublishProduct);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
