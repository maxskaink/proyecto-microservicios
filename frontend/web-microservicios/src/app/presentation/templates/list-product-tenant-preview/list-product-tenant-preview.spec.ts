import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListProductTenantPreview } from './list-product-tenant-preview';

describe('ListProductTenantPreview', () => {
  let component: ListProductTenantPreview;
  let fixture: ComponentFixture<ListProductTenantPreview>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ListProductTenantPreview]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ListProductTenantPreview);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
