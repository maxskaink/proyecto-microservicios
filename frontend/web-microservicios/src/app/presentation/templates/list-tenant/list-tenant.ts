import { ChangeDetectorRef, Component, EventEmitter, Input, Output, OnInit, OnChanges } from '@angular/core';
import { BoxTenant } from '../../components/box-tenant/box-tenant';
import { CommonModule } from '@angular/common';
import { TenantService } from '../../../service/TenantService';
import { LoadingService } from '../../../service/loading-service';
import { Tenant } from '../../../Models/Tenant';
import { FormsModule } from '@angular/forms';


@Component({
  selector: 'app-list-tenant',
  imports: [CommonModule, BoxTenant, FormsModule],
  templateUrl: './list-tenant.html',
  styleUrl: './list-tenant.css',
})
export class ListTenant implements OnChanges {
  @Output() tenantClick = new EventEmitter<string>();
  @Input() tenants: Tenant[] = [];
  
  public filteredTenants: Tenant[] = [];
  public searchQuery: string = '';

  constructor(
    private tenantService: TenantService,
    private cdr: ChangeDetectorRef,
  ) {}

  ngOnChanges(): void {
    if (this.tenants && this.tenants.length > 0) {
      this.filteredTenants = [...this.tenants];
    }
  }


  /**
   * Filtra los tenants por nombre según el término de búsqueda
   */
  onSearchChange(searchValue: string): void {
    this.searchQuery = searchValue.toLowerCase();
    this.applyFilter();
  }

  /**
   * Aplica el filtro de búsqueda
   */
  applyFilter(): void {
    if (!this.searchQuery.trim()) {
      this.filteredTenants = [...this.tenants];
    } else {
      this.filteredTenants = this.tenants.filter(tenant =>
        tenant.tenant_name.toLowerCase().includes(this.searchQuery)
      );
    }
    this.cdr.detectChanges();
  }

  /**
   * Limpia la búsqueda
   */
  clearSearch(): void {
    this.searchQuery = '';
    this.filteredTenants = [...this.tenants];
    this.cdr.detectChanges();
  }

  /**
   * Maneja el click en un tenant
   */
  handleAction(tenantId: string): void {
    this.tenantClick.emit(tenantId);
  }
}
